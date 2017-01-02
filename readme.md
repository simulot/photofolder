# Cleansing photo folder

## Introduction
After several years of shooting photos everydays with several devices, my photo repository is getting messy:
* some duplicates, same name
* some duplicates, different names
* some photo with same name, but different
* unconsistent folder naming

I wrote this utility in order to have a simple but proven folder organisation :
YEAR\YEAR.MONTH\YEAR.MONTH.DAY hierarchy based on taken date available in exif metadata when available.


Image Name collision can occurs when using several camaeras, or when the the internal counter has reached 9999.
The system should be able to  detect collision cases, and decide if this is duplicate, or an homonym.

Images are moved to the right place, and emptied folder will be removed.

## Settings
### Repository definition
   The folder naming schema has to be flexible. 

   * YYYY
   * YY 
   * MM 
   * DD 
   
   Using Go's templating system, my favorite naming schema becomes:
   /home/user/pictures/{{.YYYY}}/{{.YYYY}}.{{.MM}}/{{.YYYY}}.{{.MM}}.{{.DD}}/

### Picture name convention
   Same rule applies
   * YYYY
   * YY 
   * MM 
   * DD 
   * HH 
   * MN 
   * SS

   When several pictures have been taken during the same second, names are same.
   A number is append to the file name.

   
## Duplicate detection
  
  * Same taken date 
  * Image UniqueID when exists or MD5 hash identical 

    


## Langage choice : GO
GO language is my current favorite tool for this kind of project. Easy, portable, efficient.







