package main

import (
  "time"
  "os"
  "io/ioutil"
  "path/filepath"
  "github.com/golang/glog"
  "github.com/radovskyb/watcher"
)

type Space struct {
  scanPath        string
  watcherFiles    *watcher.Watcher
  CPU             *CPUSet
  Roles           *Roles
  Stages          *Stages
  Services        *Services
}

func NewSpace(scanPath string, enableWatcher bool) *Space {
  s := &Space{}
  s.scanPath = scanPath
  s.CPU      = NewCPUs()
  s.Roles    = NewRoles()
  s.Stages   = NewStages()
  s.Services = NewServices()
  if enableWatcher {
    s.initWatcher()
  }
  return s
}

func (s *Space) initWatcher() {
  s.watcherFiles = watcher.New()
  s.watcherFiles.SetMaxEvents(1)
  s.watcherFiles.FilterOps(watcher.Rename, watcher.Move, watcher.Remove, watcher.Create, watcher.Write)
  go func() {
    for {
      select {
      case event := <-s.watcherFiles.Event:	
        if glog.V(9) {
          glog.Infof("DBG: Watcher Event: %v", event)
        }
        s.LoadFromFiles()
      case err := <-s.watcherFiles.Error:
        glog.Fatalf("ERR: Watcher Event: %v", err)
      case <-s.watcherFiles.Closed:
        glog.Infof("LOG: Watcher Close")
        return
      }
    }
  }()
  // Start the watching process - it'll check for changes every 100ms.
  glog.Infof("LOG: Watcher Start (%s)", s.scanPath)
  if err := s.watcherFiles.AddRecursive(s.scanPath); err != nil {
    glog.Fatalf("ERR: Watcher AddRecursive: %v", err)
  }
  
  if glog.V(9) {
    // Print a list of all of the files and folders currently
    // being watched and their paths.
    for path, f := range s.watcherFiles.WatchedFiles() {
      glog.Infof("DBG: WATCH FILE: %s: %s\n", path, f.Name())
    }
  }
  go func() {
    if err := s.watcherFiles.Start(time.Millisecond * 100); err != nil {
      glog.Fatalf("ERR: Watcher Start: %v", err)
    }
  }()
}

func (s *Space) LoadFromFiles() int {
  count := 0
  errScan := filepath.Walk(s.scanPath, func(filename string, f os.FileInfo, err error) error {
    if f != nil && f.IsDir() == false  {
      ext := filepath.Ext(filename)
      if glog.V(2) {
        glog.Infof("LOG: Read file: %s", filename)
      }
      var err error
      jsonFile, err := ioutil.ReadFile(filename)
      if err != nil {
        glog.Errorf("ERR: ReadFile(%s)  #%v ", filename, err)
      } else {
        switch ext {
          case s.Services.FileExtension():
                 count += s.Services.fileParse(filename, jsonFile)
                 break
          case s.Stages.FileExtension():
                 count += s.Stages.fileParse(filename, jsonFile)
                 break
        }
      }
    }
    return nil
  })
  if glog.V(2) {
    glog.Infof("LOG: Scan Path: %s, Services: %d", s.scanPath, count)
  }
  if errScan != nil {
    glog.Errorf("ERR: ScanPath(%s): %s", s.scanPath, errScan)
  }

  return count
}
