
library(readr)
library(dplyr)
library(ggplot2)
library(tseries)

# read data 
filename <- "./data/eire_10000.0.csv"
filename <- "./data/eire_1000.0.csv"
filename <- "./data/eire_100.0.csv"
filename <- "./data/eire_10.0.csv"

filename <- "./data/gb_10.0.csv"
filename <- "./data/gb_1.0.csv"
filename <- "./data/gb_0.1.csv"
filename <- "./data/gb_0.01.csv"

filename <- "./data.csv"

walk <- read_csv(filename,
                 col_names = FALSE,
                 show_col_types = FALSE)

# tibble
names(walk) <- c('energy')
walk['time'] <- 1:nrow(walk)

# plot time series
n <- 5e04
dat <- walk[1:n,]
ggplot(data = dat,
       aes(x=time, y=energy)) +
  geom_line(col='blue')

# autocorrelation
burnin <- 100000
dat <- walk[burnin:(n+burnin),]
acf(dat$energy, lag.max=2e04)

# MC iid samples
lag <- 12000
( n <- floor((nrow(walk) - burnin)/lag) )
dat <- walk[burnin + lag*(1:n),]

# histogram samples
ggplot( data = dat,
        aes(x = energy)) +
  geom_histogram()

# mean energy
mean(dat$energy)

##################################################
# For gb_cities data we see:

# T = 10.0 (acceptance ~0.9)
# ==> burnin = 1, lag = 500, mean energy = 230.1
# T = 1.0 (acceptance ~0.3)
# ==> burnin = 1000, lag = 2000, mean energy = 119.4
# T = 0.1 (acceptance ~0.02) 
# ==> burnin = 20000, lag = 4000, mean energy = 44.9
# then gets stuck ...
# T = 0.01 ==> burnin = 20000, lag = 6000, mean energy = 45.7

# But we know we're close to the minimum energy ~ 43.0

##################################################
# For eire data we see:

# T = 10000.0 (acceptance ~0.94)
# ==> burnin = 1.2e04, lag = 2000, mean energy = 14.2e06
# T = 1000.0 (acceptance ~0.48)
# ==> burnin = 5e04, lag = 20000, mean energy = 9.6e06
# T = 100.0 (acceptance ~0.02)
# ==> burnin, lag off the scale

##################################################
# stationarity as a criterion for cooling

binsize <- 10
tolerance <- 1e-03

walk_binned <- walk %>%
  group_by(bin = floor(time/binsize)) %>%
  summarise(mu = mean(energy), sigma = sd(energy))
  
walk_binned['sigmage'] <- 0
for(i in 2:nrow(walk_binned)){
  if(walk_binned$sigma[i-1] < tolerance){
    walk_binned$sigmage[i] <- 0
  } else {
    walk_binned$sigmage[i] <- (walk_binned$mu[i] - walk_binned$mu[i-1])/walk_binned$sigma[i-1]
  }
}

ggplot(data = walk_binned[1:1000,],
       aes(x=bin, y=sigmage)) +
  geom_line(col='blue')

  






