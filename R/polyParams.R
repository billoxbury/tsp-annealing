library(readr)
library(dplyr)
library(ggplot2)
library(tseries)

# read data 
filename <- "./data/polydata_std.csv"
df_std <- read_csv(filename,
                   show_col_types = FALSE)
filename <- "./data/polydata_sigmage.csv"
df_sig <- read_csv(filename,
                   show_col_types = FALSE)

df <- df_std
  
df[df$energy == min(df$energy),] %>% View()

# histogram of polygon sizes
hist(df$npoints, breaks=20)

# histogram of energies found -
# red = 2*pi is the global minimum
# orange = 4*pi is the second harmonic (routes with winding number 2)
cond <- (df$energy < 20)
sum(cond)/nrow(df) # <--- PROPORTION SHOWN
hist(df$energy[cond], breaks=30)
rug(2*pi, col='red', lwd=10)
rug(4*pi, col='orange', lwd=10)

# running time against final energy
plot(df$time, df$energy, 
     ylim = c(0,30),
     col='blue',
     cex=0.5)
abline(2*pi, 0, col='red')
abline(4*pi, 0, col='orange')

# deviation from optimal with nr points
plot(df$npoints, df$energy, 
     ylim = c(5,30),
     cex = 0.5,
     col='blue')
abline(2*pi, 0, col='red')
abline(4*pi, 0, col='orange')


# the next 3 plots show that the result is quite robust
# to initial temperature, cooling rate and period for these problems
plot(df$cooling, df$energy, 
     ylim = c(1,100),
     cex = 0.5,
     col='blue')
plot(df$temperature, df$energy, 
     ylim = c(1,100),
     cex = 0.5,
     col='blue')
plot(df$period, df$energy, 
     ylim = c(1,100),
     cex = 0.5,
     col='blue')

# cooling rate is a combination of cooling factor and period
# let's look at final energy against these together
use <- (df$energy <= 30.0)
good <- (df$energy[use] <= 7.0)
col <- sapply(good, function(g) if(g) 'blue' else 'grey')
plot(df$period[use], df$cooling[use], 
     cex = 0.5,
     col = col)


