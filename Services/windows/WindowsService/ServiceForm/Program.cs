using Service;
using System;
using System.Diagnostics;
using System.Text.Json;

class MainClass
{
    static SocketManager sm;

    static public void Main(String[] args)
    {
        Console.WriteLine("Windows Service Started...");
        sm = new SocketManager("127.0.0.1", 8080);

        NetworkThread();

        while (true)
        {

        }
    }

    static public async Task NetworkThread()
    {
        await SignIn();

        Console.WriteLine("Started infinite listening loop...");

        while (true)
        {
            try
            {
                await ListenDataAndDoAction();
            }
            catch (Exception ex)
            {
                Console.WriteLine(ex.ToString());
            }
        }

    }

    static public async Task<bool> SignIn()
    {
        Dictionary<string, object> dict = new Dictionary<string, object>
        {
            { "username", "lavi" },
            { "password", "1111" },
            { "connectionType", "service" }
        };

        Console.WriteLine("Signing in...");
        Console.WriteLine(dict);

        var data = Methods.SerializeMessage(RequestType.LOGIN, dict);
        Console.WriteLine("Sending login... ");
        await sm.SendData(data);
        Console.WriteLine("Recieving data... ");
        var value = await sm.ReceiveData();
        Console.WriteLine("Recieved: " + value.toString());

        return true;
    }

    static public async Task<bool> ListenDataAndDoAction()
    {
        Console.WriteLine("Receiving data...");

        Data? data = await sm.ReceiveData();

        Console.WriteLine("Recieved: " + data.toString());


        // Accessing the JSON message data
        JsonDocument jsonMessage = data.GetMessageData();

        // Accessing values within the JSON
        if (jsonMessage == null)
        {
            Methods.Error("The Json Data after conversion is null!");
            return false;
        }

        if (data.code == RequestCode.Action)
        {
            string command = jsonMessage.RootElement.GetProperty("command").ToString();
            Console.WriteLine("received command: " + command);

            string[] vals = command.Split(' ');

            if (vals[0] == "start")
            {
                Console.WriteLine("Command is start, starting: " + vals[1]);
                DeviceControl.StartProgram(vals[1] + ".exe");
                return true;
            }
        }
        if (data.code == RequestCode.Query)
        {
            string query = jsonMessage.RootElement.GetProperty("query").ToString();
            Console.WriteLine("received command: " + query);

            if (query == "onBattey")
            {
                Console.WriteLine("Query is onBattery, checking battery now.");
                DeviceControl.IsRunningOnBattery();

                var response = new Data(RequestCode.QueryResponse, "{\"onBattery\": \"true\"}");
                Console.WriteLine("Sending response: "+response.toString());
                await sm.SendData(response.toBytes());

                return true;
            }
        }

        return false;
    }
}