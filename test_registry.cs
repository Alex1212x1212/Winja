using System;
using Microsoft.Win32;

class Program {
    static void Main() {
        var baseKey = Registry.CurrentUser.CreateSubKey(@"Software\Classes\*\shell\TestApp");
        baseKey.SetValue("MUIVerb", "TestApp");
        baseKey.SetValue("SubCommands", "");
        
        var shellKey = baseKey.CreateSubKey("shell");
        
        var cmd1 = shellKey.CreateSubKey("Cmd1");
        cmd1.SetValue("MUIVerb", "Command 1");
        cmd1.CreateSubKey("command").SetValue("", "notepad.exe");
        
        var cmd2 = shellKey.CreateSubKey("Cmd2");
        cmd2.SetValue("MUIVerb", "Command 2");
        cmd2.CreateSubKey("command").SetValue("", "calc.exe");
    }
}
