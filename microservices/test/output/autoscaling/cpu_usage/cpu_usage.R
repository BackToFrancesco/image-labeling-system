library(ggplot2)
library(viridis)

file_names <- c()
file_names[1] <- paste0("subtask_manager.log")
file_names[2] <- paste0("task_manager.log")

data_list <- lapply(file_names, read.csv , header=FALSE, sep=",")

plot(data_list[[1]][,1],data_list[[1]][,2],type="l",col="red", main="Overall CPU usage",
     xlab="Elapsed time (sec)", ylab="CPU usage (%)", ylim=c(0,110))
lines(data_list[[2]][,1],data_list[[2]][,2],col="green")
legend(x = "topright", 
       legend = c("Subtask manager", "Task manager"), 
       col = c("red", "green"),
       lwd = 1, 
       cex=0.8
)
