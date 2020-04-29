DONE:
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
    
TODO 
* Test on windows 10 (build 1809)  - WIP
* Build Powershell install script (https://raw.githubusercontent.com/habitat-sh/habitat/master/components/hab/install.ps1)
* Use progress (https://github.com/schollz/progressbar/issues/57)
