import os  
 
if os.name == 'posix':
    os.system('shutdown -h now')
elif os.name == 'nt':
    os.system('shutdown /s /t 0')  