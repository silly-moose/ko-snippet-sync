# KO Snippet Sync

THIS IS A WORK IN PROGRESS DO NOT USE JUST YET.

This app will sync knowledge snippets to your local desktop and then sync them back to your account.

# Golang

This app is build with GoLang version `1.14.3`. This is important for the Apple notarize process mentioned in section `Signing App For Apple`.

# Signing App For Apple

We use this cool CLI tool to Sign, notarize, and package. https://github.com/mitchellh/gon

Following the instructions on the gon page for generating a cert with Apple

After you install gon you need to create a `gon_config.json`. Here is a sample file.

```
{
    "source" : ["./ko-snippet-sync"],
    "bundle_id" : "com.example.ko-snippet-sync",
    "apple_id": {
				"provider": "ExampleCompay",
        "username" : "user@example.com",
        "password":  "XXXXXXXX"
    },
    "sign" :{
        "application_identity" : "Developer ID Application: Example, LLC"
    },
    "dmg" :{
        "output_path":  "ko-snippet-sync.dmg",
        "volume_name":  "KO Snippet Sync"
    },
    "zip" :{
        "output_path" : "ko-snippet-sync.zip"
    }
}
```
