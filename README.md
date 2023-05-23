<head>
<h1 align=center>Robin - Persistence Payload Generator</h1>
</head>

<p align="center">
  <img src="images/Robin.gif" alt="Robin"/>
</p>

Robin is a persistence enabling tool, allowing for the generation of a payload that will install a specified persistence mechanism on the target system. 

---

## Features

* Easy-to-use configuration file for payload builder
* Allows for a number of different persistence techniques
* Current supported persistence techniques:
  * Run/RunOnce Registry key
  * Winlogon UserInit Registry key
  * cmd.exe AutoRun Registry key (triggers on cmd.exe execution)
  * Service installation; ServiceDll Registry key or proxied executable
  * Startup folder LNK file

---
## Installation

Clone the repository:<br>
```git clone https://github.com/misthi0s/Robin```

Change the working directory:<br>
```cd Robin```

Build the project with Go:<br>
```go build -o Robin.exe .```

<b>NOTE:</b> The "pkgs" folder (and contents) and "config\config.yml" must exist in the same directory as the main executable. If the executable is moved, they must be moved with it.

---
## Usage

To generate a payload, simply modify the `config\config.yml` file with the appropriate settings (description of each listed below), and run the `Robin.exe` program.

---
## Configuration File Settings
Below are the settings in the configuration file and what each one controls. These descriptions also exist in the config.yml file in this repository.

<table>
<tr>
<th>Setting</th>
<th>Description</th>
<th>Options</th>
</tr>
<tr>
<td>name</td>
<td>Name of the output file when running the builder</td>
<td>N/A</td>
</tr>
<tr>
<td>technique</td>
<td>Persistence technique to install when payload is executed</td>
<td>"service", "startup", or "registry"</td>
</tr>
<tr>
<td>encryptkey</td>
<td>Password used when encrypting the configuration to pass to the payload via "-ldflags"</td>
<td>N/A</td>
</tr>
<tr>
<td>service\name</td>
<td>Name of the service to install</td>
<td>N/A</td>
</tr>
<tr>
<td>service\description</td>
<td>Description of the service that is installed</td>
<td>N/A</td>
</tr>
<tr>
<td>service\path</td>
<td>File path to the DLL or executable to run when service is started</td>
<td>N/A (* if executable is used, the execution of it will be proxied through the Robin payload (the service will be configured to execute the Robin payload on start, which will in turn run the specified executable))</td>
</tr>
<tr>
<td>service\dll</td>
<td>Specify if service payload is a DLL or executable; if DLL, the "ServiceDLL" Registry key technique will be used</td>
<td>"true" or "false"</td>
</tr>
<tr>
<td>service\runas</td>
<td>User that the service will be configured to run as</td>
<td>"SYSTEM" for LOCAL SYSTEM or "DOMAIN\User" format if service should be executed as a standard user</td>
</tr>
<tr>
<td>service\password</td>
<td>Password of standard user that is configured to run the service; will only be used if "SYSTEM" is NOT specified in "service\runas"</td>
<td>N/A</td>
</tr>
<tr>
<td>startup\name</td>
<td>Name of LNK file to be added to user's Startup folder</td>
<td>N/A</td>
</tr>
<tr>
<td>startup\path</td>
<td>File path to the executable (or command) to launch</td>
<td>N/A</td>
</tr>
<tr>
<td>registry\key</td>
<td>Registry persistence mechanism to use in the payload</td>
<td>"winlogon", "run", "runonce", or "autorun"</td>
</tr>
<tr>
<td>registry\path</td>
<td>File path to the executable to launch on persistence execution</td>
<td>N/A</td>
</tr>
<tr>
<td>registry\hive</td>
<td>Registry hive to write the Registry value to; will not be applicable to every technique</td>
<td>"HKLM" (requires administrative rights) or "HKCU"</td>
</tr>
<tr>
<td>registry\customname</td>
<td>Name of the Registry key to create, where applicable (Run/RunOnce technique)</td>
<td>N/A</td>
</tr>
</table>

---
## Future Techniques to Implement
Below are the "todo" persistence techniques to implement in future releases:

* Scheduled Tasks
* Direct Service executable
* HKCU Load Registry key
* Screensaver Registry key

---
## Additional Notes
* As mentioned in the `service\path` configuration description above, if an executable is used for a service, the Robin payload will proxy the configured executable upon service startup. This means that the service binPath will be set to wherever the Robin payload was executed from, and upon service execution, will run the Robin payload, which in turn will run the configured executable. Future releases will include the ability to create a service with a binPath pointing directly to the configured executable, but this is how it works in the current release.

---
## Issues

If you run into any issues with Robin, feel free to open an issue. Make sure to include the configuration settings used as well as any error messages produced during execution.