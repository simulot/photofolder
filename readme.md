# Cleansing photo folder

## Introduction
After several years of shooting photos everydays with several devices, my photo repository is getting messy:
* some have been copied accross directories
* some copied,with a different names
* some distinct photo have same name
* unconsistent folder naming
I wrote this utility in order to have a simple but proven 
folder organisation :
YEAR\YEAR.MONTH\YEAR.MONTH.DAY hierarchy based on taken date available in exif metadata when available.

The utility scans the given folder. For each photo, it checks the path where it should be
located. When it differs from the actual path, it moves it.

At the end of the process, empty folders are removed.

# Usage

```
usage: photofolder [<flags>] <repository> [<path>]

Flags:
      --help                     Show context-sensitive help (also try
                                 --help-long and --help-man).
  -m, --model=/{{.YYYY}}/{{.YYYY}}.{{.MM}}/{{.YYYY}}.{{.MM}}.{{.DD}}  
                                 model for path
  -d, --dryrun                   show actions to be done, but doesn't touch files
      --delete=Thumbs.db... ...  to be deleted file patterns, like thumb*.* or
                                 picasa.ini
      --delete-small             delete small image smaller than 256x256 pixels

Args:
  <repository>  media repository
  [<path>]      path to be cleaned, if empty, the whole repository is cleanned
```


## Settings
### model name definition
   * .YYYY
   * .YY 
   * .MM 
   * .DD 
   * .HH 
   * .MN 
   * .SS

   When several pictures have been taken during the same second, names are same.
   A number is append to the file name.

   
## Duplicate detection
  
  * Same taken date 
  * Image MD5 hash identical 

    


## Langage choice : GO
GO language is my current favorite tool for this kind of project. Easy, portable, efficient.







