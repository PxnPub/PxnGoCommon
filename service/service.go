package service;

import(
	OS      "os"
	Log     "log"
	Time    "time"
	Sync    "sync"
	Atomic  "sync/atomic"
	Signal  "os/signal"
	SysCall "syscall"
	ThdSafe "github.com/hayageek/threadsafe"
	Utils   "github.com/PxnPub/PxnGoCommon/utils"
);



type Service struct {
	WaitGroup   *Sync.WaitGroup
	StopChans   ThdSafe.Slice[chan bool]
	StopHooks   ThdSafe.Slice[StopHook]
	Closeables  ThdSafe.Slice[Closeable]
	CloseablesE ThdSafe.Slice[CloseableE]
	stopping    Atomic.Bool
	state       ServiceState
	timeout     int8
}

type StopHook func();


type ServiceState int8;
const (
	State_OK   ServiceState = iota
	State_Stop
	State_Warn
	State_Term
);



type AppFace interface {
	Main()
}

type Closeable interface {
	Close()
}

type CloseableE interface {
	Close() error
}



func New() *Service {
	var wait_group Sync.WaitGroup;
	service := Service{
		WaitGroup: &wait_group,
	};
	return &service;
}



func (service *Service) Start() {
	print("\n");
	go service.TrapC();
	Utils.SleepC();
}

// ctrl+c loop
func (service *Service) TrapC() {
	signals := make(chan OS.Signal, 1);
	Signal.Notify(signals, SysCall.SIGINT, SysCall.SIGTERM);
	timer := Time.NewTicker(Time.Second);
	//LOOP_TRAPC:
	for {
		SELECT_SIGNAL:
		select {
		case <-timer.C:
			service.timeout++;
			if service.timeout >= 10 {
				service.timeout = 0;
				if service.state > State_OK {
					service.state--;
				}
			}
			break SELECT_SIGNAL;
		case <-signals:
			service.state++;
			service.timeout = 0;
			SWITCH_STATE:
			switch service.state {
			case State_Stop:
				print("\r"); Log.Print("Stopping..");
				service.Stop();
				break SWITCH_STATE;
			case State_Warn:
				print("\r"); Log.Print("Terminate?");
				break SWITCH_STATE;
			default:
				if service.state < State_OK {
					service.state = State_OK;
				} else
				if service.state > State_Warn {
					Utils.SleepC()
					print("\r"); Log.Print("Terminated!!!");
					OS.Exit(0);
				}
				break SWITCH_STATE;
			}
			break SELECT_SIGNAL;
		}
	} // end LOOP_TRAPC
}

func (service *Service) Stop() {
	service.stopping.Store(true);
	var finished bool;
	LOOP_STOP:
	for {
		finished = true;
		// chan
		LOOP_STOPCHANS:
		for ; service.StopChans.Length()>0; {
			stopchan, ok := service.StopChans.Get(0);
			if !ok { break LOOP_STOPCHANS; }
			finished = false;
			service.StopChans.Remove(0);
			stopchan <-true;
		}
		// Close()
		LOOP_CLOSEABLES:
		for ; service.Closeables.Length()>0; {
			closeable, ok := service.Closeables.Get(0);
			if !ok { break LOOP_CLOSEABLES; }
			finished = false;
			service.Closeables.Remove(0);
			closeable.Close();
		}
		// Close() error
		LOOP_CLOSEABLES_E:
		for ; service.CloseablesE.Length()>0; {
			closeable, ok := service.CloseablesE.Get(0);
			if !ok { break LOOP_CLOSEABLES_E; }
			finished = false;
			service.CloseablesE.Remove(0);
			closeable.Close();
		}
		// func()
		LOOP_STOPHOOKS:
		for ; service.StopHooks.Length()>0; {
			hook, ok := service.StopHooks.Get(0);
			if !ok { break LOOP_STOPHOOKS; }
			finished = false;
			service.StopHooks.Remove(0);
			hook();
		}
		if finished { break LOOP_STOP; }
		Utils.SleepX();
	}
}

func (service *Service) IsStopping() bool {
	return service.stopping.Load();
}



func (service *Service) Wait() {
	service.WaitGroup.Wait();
}

func (service *Service) WaitUntilEnd() {
	Utils.SleepC(); print("\n"); service.Wait();
	Utils.SleepC(); print(" ~end~ \n");
	print("\n"); OS.Exit(0);
}



func (service *Service) NewStopChan() chan bool {
	stopchan := make(chan bool, 1);
	if service.stopping.Load() { stopchan <-true; }
	service.StopChans.Append(stopchan);
	return stopchan;
}

func (service *Service) AddStopHook(hook StopHook) {
	if service.stopping.Load() {
		hook();
	} else {
		service.StopHooks.Append(hook);
	}
}

func (service *Service) AddClose(closeable Closeable) {
	if service.stopping.Load() { closeable.Close();
	} else { service.Closeables.Append(closeable); }
}

func (service *Service) AddCloseE(closeable CloseableE) {
	if service.stopping.Load() {
		if err := closeable.Close(); err != nil {
			Log.Printf("%s", err); }
	} else {
		service.CloseablesE.Append(closeable);
	}
}
