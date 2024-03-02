library(ggplot2)
library(viridis)

file_names <- c()

for (i in 1:10) {
  file_names[i] <- paste0("classifier-", i-1, ".log")
}

data_list <- lapply(file_names, read.csv , header=FALSE, sep=",")

plot(data_list[[1]][,2],data_list[[1]][,1],type="l",col="red", main="Operator Response Times",
     xlab="Elapsed time (sec)", ylab="Response time (sec)", ylim=c(0,5))
lines(data_list[[2]][,2],data_list[[2]][,1],col="green")
lines(data_list[[3]][,2],data_list[[3]][,1],col="yellow")
lines(data_list[[4]][,2],data_list[[4]][,1],col="black")
lines(data_list[[5]][,2],data_list[[5]][,1],col="blue")
lines(data_list[[6]][,2],data_list[[6]][,1],col="purple")
lines(data_list[[7]][,2],data_list[[7]][,1],col="pink")
lines(data_list[[8]][,2],data_list[[8]][,1],col="orange")
lines(data_list[[9]][,2],data_list[[9]][,1],col="lightblue")
lines(data_list[[10]][,2],data_list[[10]][,1],col="gray")
legend(x = "topright", 
       legend = c("classifier-1", "classifier-2", "classifier-3", "classifier-4", "classifier-5", "classifier-6", "classifier-7", "classifier-8", "classifier-9", "classifier-10"), 
       col = c("red", "green", "yellow", "black", "blue", "purple", "brown", "pink", "orange", "lightblue", "gray"),
       lwd = 1, 
       cex=0.8
)