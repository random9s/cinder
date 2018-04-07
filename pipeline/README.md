### Pipeline
> Simple API that allows easy usage and control of multiple channels between processes in a pipeline

#### Structs:

* Pipeline - Responsible for managing all channels and processes
* Process - Responsible for processing arbitrary data and passing data to the next section of the pipeline

#### Pipeline API:

* New *Pipeline
> Initializes a new pipeline
> * Side Effects: 
> - Creates a nil process to use as the head of the pipeline
> - Starts a goroutine which waits for an error to occur
> * NOTE: This is the only way a pipeline should be initialized

* Start 
> Uses the initialized nil process to send data to the first, user created, process
> * Side Effects:
> - Attaches a "closer" channel to the tail process(es), this only happens on the first call of start

* Wait error 
> Wait closes the initialized nil process and waits for all tail processes to send the closer signal.  Once all procs have closed, we send back an error, if one occurred
> * NOTE: Start will panic if called after wait

* WaitWithTimeout
> This is the same as wait but with an added timeout

* Abort
> Abort sends an error to the pipelines abort channel.  This causes the pipeline to prematurely close all running processes and returns the error from the Wait func
> * NOTE: Access to this function is provided to all the process functions

### Creating a custom pipeline:
* Append error
> This function will append a custom process to the pipelines internal process slice
> This provides a simple 1 to 1 connection between two processes

* FanOut error
> This function will append N custom processes to the previous process in the pipeline
> This provides a 1 to N connection between multiple processes

* FanIn error
> This function will append 1 process to the previous N processes in the pipeline
> This provides an N to 1 connection between multiple processes

* ConnectNtoM error
> This function will append M processes to the previous N processes in the pipeline
> This, as its name says, provides an M to N connection between multiple processes


#### Creating a custom process:

* NewProcess *Process
> Initializes and returns a custom process with the appropriate Process function
> All you need to create a custom process is to wrap your function with a process function
> A Process function follows the form: func(v interface{}, send func(interface)bool, abort func(error))
> v is the data that is passed from process to process
> send contains safe access the processes send channel, returns true if the channel has been closed
> abort contains safe access to the pipelines abort channel
