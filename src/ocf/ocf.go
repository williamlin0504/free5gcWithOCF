package main

import (
    "fmt"
    "html"
    "log"
    "net/http"
    "free5gc/src/app"
) 

func main() {
    //Mount OCF to free5gc
    app := cli.NewApp()
	app.Name = "ocf"
	fmt.Print(app.Name, "\n")
	appLog.Infoln("OCF version: ", version.GetVersion())
	app.Usage = "-free5gccfg common configuration file -ocfcfg ocf configuration file"
	app.Action = action
	app.Flags = OCF.GetCliCmd()

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("OCF Run err: %v", err)
	}
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
    })

    log.Println("OCF Server Started. Listening on localhost: 8080")
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func Nchf_ConvergedChargingFunction_create(ue_ID){
    resp, err := http.PostForm("https://je752rauad.execute-api.us-east-1.amazonaws.com/Nchf/create",url.Values{"key": {"ue-ID"}, "id": {ue_ID}})

    if err != nil {
        fmt.Print(err.Error()) 
        os.Exit(1)
    }

    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }

    if responseData == 1{
        log.Println("Request Authorized.")
        response, err := http.PostForm("https://je752rauad.execute-api.us-east-1.amazonaws.com/Nchf/continous-write",url.Values{"key": {"ue-ID"}, "id": {ue_ID}})

    if err != nil {
        fmt.Print(err.Error()) 
        os.Exit(1)
    }

    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }

    else if responseData == 0{
        log.Println("User does not have enough GU.")
        notifyAPI()
    }
}

func Nchf_ConvergedChargingFunction_update(ue_ID){
    //User tends to update back to GU 50
    response, err := http.PostForm("https://je752rauad.execute-api.us-east-1.amazonaws.com/Nchf/update",url.Values{"key": {"ue-ID"}, "id": {ue_ID}})

    if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }

    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }
    log.Println(string(responseData))
}

func Nchf_ConvergedChargingFunction_release(ue_ID){
    //When user stops the session
    response, err := http.POST("https://je752rauad.execute-api.us-east-1.amazonaws.com/Nchf/release",url.Values{"key": {"ue-ID"}, "id": {ue_ID}})

    if err != nil { 
        fmt.Print(err.Error())
        os.Exit(1)
    }

    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Session Released!")
}

func Nchf_ConvergedChargingFunction_notify(ue_ID){
    //USELESS TIL I FIGURE OUT
}