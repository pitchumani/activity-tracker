# Interview Preparation - Linux

## systemd
systemd is the init system and system manager used in many Linux distributions. It is the first user-space process started by the kernel during boot (PID 1) and is responsible for initializing and managing system services, daemons, and other system components. systemd aims to provide a more efficient and robust way to manage the Linux system boot process and runtime environment compared to older init systems like SysVinit.
> Init system & service manager. Boots services in parallel, manages dependencies, logs via journald.

## systemctl
systemctl is the command-line utility used to interact with and control the systemd init system. It allows users and administrators to manage systemd "units," which are configuration files defining how systemd handles various resources, such as services, mount points, sockets, and targets.
### Example usages

* Checking the status of a service: ```systemctl status nginx.service```. Displays detailed information about nginx web server service, including its current status (active, inactive, failed), whether it's enabled to start at boot, and recent log messages.

* Staring a service:```sudo systemctl start apache2.service```. This command initiates the `apache2` web server service.

* Stopping a service: ```sudo systemctl stop ssh.service```. This command terminates the `ssh` daemon service.

* Enable a service to start at boot. ```sudo systemctl enable firewalld.service```. This command configures the `firewalld` service to automatically start whenever the system boots.

* Disabling a service from starting at boot: ```sudo systemctl disable cups.service```. This command prevents the cups printing service from automatically starting at boot.

* Restarting a service. ```sudo systemctl restart mysql.service```.
This command stops and then restarts the mysql database service.

* Listing all loaded units. ```systemctl list-units```.
This command displays a list of all currently loaded systemd units, including services, mounts, and targets, along with their status.

### fork() vs exec()
#### fork():
This system call creates a new process, known as the child process, which is an almost exact duplicate of the calling process (the parent). The child process inherits a copy of the parent's memory space, open file descriptors, and other attributes. After fork(), both the parent and child processes continue execution independently, often with the child process executing different code or performing specific tasks.

#### exec():
This family of system calls (e.g., execve, execlp) replaces the current process image with a new executable program. When an exec() call is successful, the code, data, and stack of the calling process are entirely overwritten by the new program. The Process ID (PID) of the process remains the same, but the program being executed changes.

#### Key Differences:
**Process Creation vs. Replacement:**
fork() creates a new process, while exec() replaces the current process with a different program.

**PID:** fork() results in a new, distinct PID for the child process, whereas exec() maintains the same PID for the process that is being replaced.

**Execution Flow:**
After fork(), both parent and child processes continue executing. After exec(), the original program's execution ceases, and the new program begins execution within the same process.

**Common Usage:**
The fork() and exec() calls are frequently used together, particularly in shell environments. A common pattern involves a parent process using fork() to create a child process, and then the child process using exec() to load and execute a new program (e.g., when a user types a command in a shell, the shell `fork()`s and the child `exec()`s the command). This allows the parent process (the shell) to continue running while the new program executes in its own process space.

#### What happens when a user types a command in a shell
When a user types a command in a shell and presses Enter, a series of steps are initiated by the shell:

**Input Reading and Tokenization:**
The shell first reads the entire line of input entered by the user. It then performs tokenization, breaking down the command line into individual components or "tokens," such as the command itself, arguments, and operators.

**Parsing and Interpretation:**
The shell parses the tokens, interpreting their meaning and identifying any special characters or shell-specific constructs (like variables, wildcards, or command substitutions).

**Command Lookup:**
The shell attempts to locate the command. It checks for:
* Aliases: User-defined shortcuts for commands.
* Functions: Shell functions defined by the user or system.
* Built-in commands: Commands that are part of the shell itself (e.g., cd, echo).
* Executable files in the PATH: If not found as an alias, function, or built-in, the shell searches for an executable file with the given name in the directories specified by the PATH environment variable.

#### Explain inode and how file permissions work
An inode (index node) is a data structure in Unix-like file systems that stores metadata about a file or directory. Each file and directory has a unique inode number. This metadata includes:
* File type: Indicates if it's a regular file, directory, symbolic link, etc.
* Permissions: Defines who can read, write, or execute the file.
* Ownership: Specifies the user and group that own the file.
* Timestamps: Records creation, modification, and last access times.
* File size: The size of the file in bytes.
* Number of hard links: The count of directory entries pointing to this inode.
* Pointers to data blocks: Locations on the disk where the actual file content is stored.

**File Permissions:**
File permissions in Unix-like systems control access to files and directories. They are assigned to three categories of users:
* Owner: The user who owns the file.
* Group: Users belonging to the group associated with the file.
* Others: All other users on the system.

For each category, three types of permissions can be granted:
* Read (r): Allows viewing the file's content or listing a directory's contents.
* Write (w): Allows modifying the file's content or creating/deleting files within a directory.
* Execute (x): Allows running an executable file or entering a directory.

Permissions are often represented in a symbolic format (e.g., rwx, rw-, r--) or as an octal number (e.g., 7 for rwx, 6 for rw-, 4 for r--). For example, chmod 755 filename grants read, write, and execute permissions to the owner, and read and execute permissions to the group and others.

#### Difference between hard link and soft link
Hard links and soft links (also known as symbolic links) are both ways to create aliases for files or directories in a file system, but they differ significantly in how they operate. A hard link creates a new directory entry that points directly to the same inode (data block) as the original file, effectively making it another name for the same data. A soft link, on the other hand, acts as a pointer to the original file's path, like a shortcut. 

Here's a more detailed breakdown: 
##### Hard Links: 
**Direct Data Reference:**
A hard link points directly to the file's inode, which contains the file's metadata and data location. This means that the hard link and the original file share the same inode. 
**Same Inode, Same File:**
Modifying a file through a hard link will affect the original file and vice-versa. They are essentially the same file, just with different names.

**File System Bound:**
Hard links cannot cross file system boundaries (e.g., from one hard drive partition to another). 

**Cannot Link Directories:**
Hard links cannot be created for directories due to the potential for creating circular dependencies. 

**Space Efficiency:**
Hard links don't create a new copy of the data, so they consume minimal extra space (only the directory entry). 

**Deletion:**
Deleting the original file or any hard link pointing to it only removes the link from the directory entry. The data remains on the disk until all hard links are removed. 

##### Soft Links (Symbolic Links):
**Path Reference:** A soft link stores the path to the original file or directory. It's essentially a shortcut. 

**Independent Entity:** Soft links are separate files. Changes made through a soft link don't directly affect the original file (unless the soft link points to a file that is subsequently modified directly). 

**Cross File System:** Soft links can span file system boundaries. 

**Can Link Directories:** Soft links can link to both files and directories. 

**Space Overhead:** Soft links consume a small amount of space to store the path, but not as much as a full copy of the file. 

**Vulnerability:** If the original file or directory is moved or deleted, the soft link becomes broken (dangling) and will no longer work. 

##### In essence: 
Hard links are like having multiple names for the same file or data on the same disk, and all names behave identically. 
Soft links are like shortcuts that point to the location of a file or directory. 

### What is ldd and when to use it?
`ldd` is a command-line utility on Unix-like systems that lists the shared libraries required by a program or another shared library. It's primarily used for troubleshooting missing dependencies or verifying library compatibility during deployment. 

Here's a more detailed explanation:
**What it does:**
* ldd essentially shows you which shared object files (libraries) a program depends on to run correctly. 
* It does this by invoking the dynamic linker with a special environment variable, which then lists the shared libraries that the linker would load. 
* ldd doesn't execute the program itself, so it won't show libraries loaded dynamically at runtime using dlopen(). 

**When to use it:**
> Troubleshooting "missing library" errors:

When a program fails to start with an error like "error while loading shared libraries: libfoo.so.1: cannot open shared object file," ldd can help identify the missing library. 

> Checking for unnecessary dependencies:

The -u option helps identify libraries that are linked but not actually used, potentially leading to cleaner deployments. 

> Preparing application deployment:

ldd can be used to create a list of required libraries for an application, ensuring all dependencies are available on the target system. 

> Verifying library compatibility:

The -v option shows version information, which can be useful for confirming that the correct library versions are installed. 

> Security considerations:

While ldd is a useful tool, it's crucial to avoid using it on untrusted executables. Some implementations of ldd may execute the program, potentially leading to the execution of arbitrary code. If you need to check dependencies of an untrusted program, objdump -p is a safer alternative, according to Stack Overflow article.

### netstat
netstat (network statistics) is a command-line utility used to display network connections (both incoming and outgoing), routing tables, interface statistics, and multicast memberships on a system. It is available on various operating systems, including Windows, Linux, and macOS. It is a valuable tool for network troubleshooting, performance monitoring, and security auditing.
* all active connections and listening ports: `netstat -a`
* stats for specific protocols `netstat -s p tcp`
