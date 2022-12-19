#define MyAppName "IMLS Session Counter"
#define MyAppVersion "1.0"
#define MyAppPublisher "GSA 10x"
#define MyAppURL "https://github.com/IMLS/estimating-wifi"
#define MyAppExeName "windows-session-counter.exe"
#define MySecondaryAppExeName "wifi-hardware-search-windows.exe"

[Setup]
AppId={{8D2CDEA5-9C55-44D4-84B3-ACDE9D4035BD}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
;AppVerName={#MyAppName} {#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}
DefaultDirName={autopf}\{#MyAppName}
DisableProgramGroupPage=yes
;PrivilegesRequired=lowest
OutputBaseFilename=SessionCounterInstall
OutputDir=.
Compression=lzma
SolidCompression=yes
WizardStyle=modern
SetupLogging=yes

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Dirs]
Name: "{app}\Wireshark"
Name: "{app}\service"

[Files]
; NOTE: Do not use "Flags: ignoreversion" on any shared system files
; Our IMLS installer
Source: "{#MyAppExeName}"; \
  DestDir: "{app}"; \
  Flags: ignoreversion
Source: "{#MySecondaryAppExeName}"; \
  DestDir: "{app}"; \
  Flags: ignoreversion
Source: "README.md"; \
  DestDir: "{app}"; \
  Flags: ignoreversion
Source: "session-counter.ini"; \
  DestDir: "{app}\service"; \
  Flags: ignoreversion; \
  AfterInstall: WriteOutIni
Source:"WinSw-x64.exe"; \
  DestDir: "{app}\service"; \
  Flags: ignoreversion
Source:"WinSw-x64.xml"; \
  DestDir: "{app}\service"; \
  Flags: ignoreversion
; Wireshark 3.6.5 portable app
Source:"Wireshark-win64-4.0.0.exe"; \
  DestDir: "{app}\Wireshark"; \
  Flags: ignoreversion

[Run]
;Filename: "{app}\{#MySecondaryAppExeName}"; \
  ;Description: "wifi-hardware-search-windows"; \
  ;Flags: runascurrentuser
Filename: "{app}\Wireshark\Wireshark-win64-4.0.0.exe"; \
  Description: "Wireshark 4.0.0"; \
  Parameters: "\D C:\imls"; \
  Flags: runascurrentuser
Filename: "{app}\service\WinSw-x64.exe"; \
  Parameters: "install"; \
  Description: "WinSw-x64 install"; \
  Flags: runascurrentuser
Filename: "{app}\service\WinSw-x64.exe"; \
  Parameters: "start"; \
  Description: "WinSw-x64 start"; \
  Flags: runascurrentuser

[UninstallDelete]
Type: files; Name: "{app}\service\WinSw-x64.wrapper.log"
Type: files; Name: "{app}\service\WinSw-x64.err.log"
Type: files; Name: "{app}\service\WinSw-x64.out.log"

[UninstallRun]
Filename: "{app}\service\WinSw-x64.exe"; \
  Parameters: "stop --force"; \
  Flags: runascurrentuser
Filename: "{app}\service\WinSw-x64.exe"; \
  Parameters: "dev kill"; \
  Flags: runascurrentuser
Filename: "{app}\service\WinSw-x64.exe"; \
  Parameters: "uninstall"; \
  Flags: runascurrentuser

[Code]

var
  IntroPage: TOutputMsgWizardPage;
  LibraryPage: TInputQueryWizardPage;

procedure InitializeWizard;
begin
  IntroPage := CreateOutputMsgPage(wpWelcome,
    'Session Counter Pilot Program',
    'Thank you for participating in the Session Counter pilot!',
    'Before proceeding with installation, please confirm that you have gotten ' +
    'an email from your state library director which contains your API key. ' +
    'If not, please stop the installation by clicking ''Cancel''.');

  LibraryPage := CreateInputQueryPage(IntroPage.ID,
    'Library Information',
    'This information will help uniquely identify your library in the state system.',
    'Please enter your API key and your public library Federal-State Cooperative ' +
    'System (FSCS) ID, then click Next.');
  LibraryPage.Add('API key:', False);
  LibraryPage.Add('FSCS ID:', False);
end;

function IsAlphabetic(C: Char): Boolean;
begin
   { InnoSetup does not support char ranges. }
   if (C >= 'a') and (C <= 'z') then begin
      Result := True;
   end else if (C >= 'A') and (C <= 'Z') then begin
      Result := True;
   end else begin
      Result := False;
   end;
end;

function IsNumeric(C: Char): Boolean;
begin
   { InnoSetup does not support char ranges. }
   if (C >= '0') and (C <= '9') then begin
      Result := True;
   end else begin
      Result := False;
   end;
end;

function NextButtonClick(CurPageID: Integer): Boolean;
var
   Temp: String;
begin
  if CurPageID = LibraryPage.ID then begin
    { Check for empty data }
    if Trim(LibraryPage.Values[0]) = '' then begin
      MsgBox('You must enter an API key.', mbError, MB_OK);
      Result := False;
    end else begin
      if Trim(LibraryPage.Values[1]) = '' then begin
        MsgBox('You must enter an FSCS ID.', mbError, MB_OK);
        Result := False;
      end else begin
        { Check for formatting }
        Temp := Trim(LibraryPage.Values[1]);
        if Length(Temp) <> 10 then begin
           MsgBox('You must enter a FSCS ID in the format AB1234-567. ' +
                  'Error: ID is too short', mbError, MB_OK);
           Result := False;
        end else begin
          if not (IsAlphabetic(Temp[1]) and
                  IsAlphabetic(Temp[2])) then begin
             MsgBox('You must enter a FSCS ID in the format AB1234-567. ' +
                    'Error: ID does not have a state', mbError, MB_OK);
             Result := False;
          end else if not (IsNumeric(Temp[3]) and
                  IsNumeric(Temp[4]) and
                  IsNumeric(Temp[5]) and
                  IsNumeric(Temp[6])) then begin
             MsgBox('You must enter a FSCS ID in the format AB1234-567. ' +
                    'Error: ID does not have a 4 digit number', mbError, MB_OK);
             Result := False;
          end else if (Temp[7] <> '-') then begin
             MsgBox('You must enter a FSCS ID in the format AB1234-567. ' +
                    'Error: ID does not have a dash (-)', mbError, MB_OK);
             Result := False;
          end else if not (IsNumeric(Temp[8]) and
                  IsNumeric(Temp[9]) and
                  IsNumeric(Temp[10])) then begin
             MsgBox('You must enter a FSCS ID in the format AB1234-567. ' +
                    'Error: ID does not have a 3 digit number', mbError, MB_OK);
             Result := False;
          end else begin
             Result := True;
          end;
        end;
      end;
    end;
    end else
      Result := True;
end;

procedure WriteOutIni();
begin
  SetIniString('device', 'api_key', LibraryPage.Values[0], ExpandConstant(CurrentFileName));
  SetIniString('device', 'fscs_id', LibraryPage.Values[1], ExpandConstant(CurrentFileName));
end;

var
  ExitCode: Integer;
  FullFilePath: String;

procedure CurStepChanged(CurStep: TSetupStep);
begin
 if CurStep = ssInstall then
  begin
    FullFilePath:= GetCurrentDir + '\wifi-hardware-search-windows.exe';
    //Run wifi-hardware-search
    if ExecAsOriginalUser(
      ExpandConstant(FullFilePath), '', '', SW_SHOW, ewWaitUntilTerminated, ExitCode) then
    begin
      if ExitCode <> 0 then begin
        //wifi-hardware-search failed, abort install
        SuppressibleMsgBox('Failed to find hardware device, aborting install.', mbError, MB_OK, IDOK);
        Abort;
      end;
    end;
  end;
end;
