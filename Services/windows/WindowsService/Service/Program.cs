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
    }

    static public async void NetworkThread()
    {
        await SignIn();

        while (true)
        {
            await ListenDataAndDoAction();
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
        await sm.SendData(data);
        await sm.ReceiveData();

        return true;
    }

    static public async Task<bool> ListenDataAndDoAction()
    {
        Data? data = await sm.ReceiveData();

        if (data == null) return false;

        Console.WriteLine(data);

        // Accessing the JSON message data
        JsonDocument jsonMessage = data.GetMessageData();

        // Accessing values within the JSON
        if (jsonMessage == null)
        {
            Methods.Error("The Json Data after conversion is null!");
            return false;
        }
        string command = jsonMessage.RootElement.GetProperty("command").GetString()!;

        Console.WriteLine("recieved command: "+command);

        string[] vals = command.Split(' ');
        if (vals[0] == "start")
        {
            Console.WriteLine("Command is start, starting: "+ vals[1]);
            ProgramStarter.StartProgram(vals[1] + ".exe");
            return true;
        }

        return false;
    }    
}