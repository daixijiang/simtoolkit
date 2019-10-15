# simtoolkit
sim produce tool (daixijiang@gmail.com)  
https://github.com/nskygit/simtoolkit

![image](https://github.com/nskygit/simtoolkit/raw/master/simtoolkit.png)

## Requirements
* **GO**
* **mingw**
* **GCC + Make** - GNU C Compiler and build automation tool
* **simcrypt.dll (private)**

## Usage
**1 start simtoolkit**

    cmd>simtoolkit.exe
    run simtoolkit.exe

**2 run simtoolkit**

    open serial; produce; check  

**3 config (simconfig.toml)**

    verbose            = 1
    #simfake            = 1
    module             = "ec20"

    [serial]
    serial_max         = 8
    serial_timeout     = 3000
    serial_timewait    = 200

    [produce]
    timeout_cold_reset = 0
    timeout_hot_reset  = 0
    timeout_creg       = 3
    timeout_common     = 1
