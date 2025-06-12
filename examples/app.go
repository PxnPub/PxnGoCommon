package app;

import(
	Service "github.com/PxnPub/PxnGoCommon/service"
);



type MyApp struct {
}



func New() Service.App {
	return &MyApp{};
}

func (app *MyApp) Main() {
	service := Service.New();
	service.Start();

print("test works!\n");

	service.WaitUntilEnd();
}
