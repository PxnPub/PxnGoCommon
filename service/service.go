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
	WaitGroup *Sync.WaitGroup
	StopChans ThdSafe.Slice[chan bool]
	StopHooks ThdSafe.Slice[StopHook]
	Stopping  Atomic.Bool
	State     ServiceState
	Timeout   int8
}

type StopHook func();


type ServiceState int8;
const (
	State_OK   ServiceState = iota
	State_Stop
	State_Warn
	State_Term
);



type App interface{
	Main()
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
			service.Timeout++;
			if service.Timeout >= 10 {
				service.Timeout = 0;
				if service.State > State_OK {
					service.State--;
				}
			}
			break SELECT_SIGNAL;
		case <-signals:
			service.State++;
			service.Timeout = 0;
			SWITCH_STATE:
			switch service.State {
			case State_Stop:
				print("\r"); Log.Print("Stopping..");
				service.Stop();
				break SWITCH_STATE;
			case State_Warn:
				print("\r"); Log.Print("Terminate?");
				break SWITCH_STATE;
			default:
				if service.State < State_OK {
					service.State = State_OK;
				} else
				if service.State > State_Warn {
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
	service.Stopping.Store(true);
	for ; service.StopHooks.Length()>0; {
		stopchan, ok := service.StopChans.Get(0);
		if !ok { break; }
		service.StopChans.Remove(0);
		stopchan <-true;
	}
	for ; service.StopHooks.Length()>0; {
		hook, ok := service.StopHooks.Get(0);
		if !ok { break; }
		service.StopHooks.Remove(0);
		hook();
	}
}

func (service *Service) IsStopping() bool {
	return service.Stopping.Load();
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
	if service.Stopping.Load() {
		stopchan <-true;
	}
	service.StopChans.Append(stopchan);
	return stopchan;
}

func (service *Service) AddStopHook(hook StopHook) {
	if service.Stopping.Load() {
		hook();
	} else {
		service.StopHooks.Append(hook);
	}
}
