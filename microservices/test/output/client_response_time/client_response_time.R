library(ggplot2)
library(viridis)

file_names <- c()

for (i in 1:15) {
  file_names[i] <- paste0("client-", i-1, ".log")
}

data_list <- lapply(file_names, read.csv , header=FALSE, sep=",")

plot(data_list[[1]][,2],data_list[[1]][,1],type="l",col="red", main="Client Response Times",
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
lines(data_list[[11]][,2],data_list[[11]][,1],col="cyan")
lines(data_list[[12]][,2],data_list[[12]][,1],col="darksalmon")
lines(data_list[[13]][,2],data_list[[13]][,1],col="coral")
lines(data_list[[14]][,2],data_list[[14]][,1],col="deeppink")
lines(data_list[[15]][,2],data_list[[15]][,1],col="darkseagreen1")
legend(x = "topright", 
       legend = c("client-1", "client-2", "client-3", "client-4", "client-5", "client-6", "client-7", "client-8", "client-9", "client-10", "client-11", "client-12", "client-13", "client-14", "client-15"), 
       col = c("red", "green", "yellow", "black", "blue", "purple", "brown", "pink", "orange", "lightblue", "gray", "cyan", "darksalmon", "coral", "deeppink", "darkseagreen1"),
       lwd = 1, 
       cex=0.8
)