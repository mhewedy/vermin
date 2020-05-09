## TODO: 
* Write more test cases
* Work on Clone to clone a VM (export: [vboxmanage export vm_01 --ovf20 -o ~/Documents/temp.ova] then import)


## DONE:
*  Build CI
*  Build Linux/mac install script 
*  Eliminate need to command dependecies (most of them):
    * arp   (exists on win/nix)
    * ping  (exists on win/nix)
    * vboxmanage (a dependency)
    * wget      (exists on win/nix) - on windows (wget <file url> -o <file output>)
    * ssh		(exists on win/nix)
    * scp       (exists on win/nix)
*  Logo drawing
* Fix issue of invalid ssh session when login after start
* Build PowerShell install script (https://raw.githubusercontent.com/habitat-sh/habitat/master/components/hab/install.ps1)
* Test on windows 10 (build 1809)  - (Partially Done)
* Enable mounting:
    on guest:
    ```
    sudo apt-get install virtualbox-guest-utils
    sudo adduser $USER vboxsf
    ```
    on host:
    ```
    vboxmanage sharedfolder add vm_01 --name <unique_name_e.g._pwd> --hostpath $(pwd) --transient --automount --auto-mount-point /vermin
    # --transient to be an option
    ```
* Use progress (https://github.com/schollz/progressbar/issues/57)
* Depend on *.vbox file for reading all information related to VM (MAC, etc...)    

## POSTPONED:
* consider use google drive (https://drive.google.com/uc?export=download&confirm=htAy&id=<fileid>) 
