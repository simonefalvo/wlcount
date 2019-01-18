# wlcount

Word length count distributed RPC application.

## Getting Started

### Prerequisites
* [Go](https://golang.org/)

### Installing

Download/clone the repository and build the [wlcount](wlcount.go) program
```
go build wlcount.go
```
Optionally you can build the [worker](worker.go) program too (the script [launch_workers.sh](launch_workers.sh) will do it as well)

```
go build worker.go
```

## Running

To run the application first launch a set of workers using the shell script [launch_workers.sh](launch_workers.sh) specifying the number of workers you want to deploy
```
./lauch_workers.sh 5
```
With the above code example, 5 workers will be launched on the localhost and will be listening on OS-assigned ports.
The wlcount client will read a local configuration file (automatically generated) in order to get the addresses and connect to them, to run it simply type the following code on a different terminal (in order to not confuse programs' printings)
```
./wlcount file1 [file2 ...]
```
where *file1* and *file2* are the names of the files you want to count the words' lengths.
Otherwise the workers can be run manually either on several terminals or in the same one running them in background, but it is important to remember to **delete *address.config*** file before each time you want to run a set of workers.

### Running Example
First launch a set of 5 workers, they will be sent in backround.
```
~/go/src/github.com/smvfal/wlcount$ ./launch_workers.sh 5
shutting down previous workers...
worker: no process found
launching 5 workers...
done
22864 pts/2    00:00:00 launch_workers.
22867 pts/2    00:00:00 worker
22869 pts/2    00:00:00 worker
22876 pts/2    00:00:00 worker
22890 pts/2    00:00:00 worker
22892 pts/2    00:00:00 worker
```
Then run the main program specifying the example file ([*examples/short_meta.txt*](examples/short_meta.txt) in this case)
```
~/go/src/github.com/smvfal/wlcount$ ./wlcount examples/short_meta.txt 
getting workers..
workers: [[::]:43823 [::]:44201 [::]:45515 [::]:45971 [::]:33807]
Counting file examples/short_meta.txt
Mapping...
Sorted word lengths: [1 2 3 4 5 6 7 8 9 10 11]
Reducing...
---------------
Length | Count 
-------+-------
     1 | 2
     2 | 15
     3 | 19
     4 | 14
     5 | 12
     6 | 10
     7 | 4
     8 | 5
     9 | 1
    10 | 1
    11 | 2
---------------
```
And on the server(worker) side you should have an output like this:
```
22867 :  mapping..
22869 :  mapping..
22892 :  mapping..
22867 :  map complete!
22869 :  map complete!
22892 :  map complete!
22890 :  mapping..
22890 :  map complete!
22876 :  mapping..
22876 :  map complete!
22867 :  reducing..
22867 :  reduce complete!
22869 :  reducing..
22869 :  reduce complete!
22876 :  reducing..
22876 :  reduce complete!
22890 :  reducing..
22892 :  reducing..
22890 :  reduce complete!
22892 :  reduce complete!
```
The number on the left side is the PID of the worker executing the task.
## Authors

* **Simone Falvo**

## License

This project is licensed under the GNU GPL v3.0 License - see the [LICENSE](LICENSE) file for details

## Acknowledgments

* [MapReduce: Simplified Data Processing on Large Clusters - *Jeffrey Dean and Sanjay Ghemawat*](https://static.googleusercontent.com/media/research.google.com/it//archive/mapreduce-osdi04.pdf)
