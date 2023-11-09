
library(readr)
library(dplyr)
library(ggplot2)

# read data 
#filename <- "./data/gb_1.0.csv"
filename <- "./data/10k-gon-T0pt1.csv"
df <- read_csv(filename,
                 show_col_types = FALSE)
df['walker'] <- factor(df$walker)

# plot by temperature - but only m~12 values
temp_set <- rev(sort(unique(df$temperature)))
n <- length(temp_set)
m <- 10
temp_subset <- temp_set[seq(1,n,ceiling(n/m))]

# period (after sampling)
period <- 1 + max(df['iteration'])

# partition to burn-in and mature data frames
df_burnin <- df %>%
  filter( iteration < period/2 )
df_mature <- df %>%
  filter( iteration > period/2 )

# boxplot comparison of walkers conditioned on temperature
df_mature %>% 
  filter(temperature %in% temp_subset) %>%
  ggplot(aes(group=walker, energy)) +
  geom_boxplot() +
  facet_wrap(~temperature) + 
  theme(
    axis.text.y = element_blank(),
    axis.ticks.y = element_blank())

# function to extract F-value and p-value from aov() output
# (to work out this function just examine unlist(summary()))
fv <- function(aov_output){
  
  unlist(summary(aov_output))[c(7,9)]
  #as.numeric( unlist(summary(aov_output))[c(7,9)] )
  
}

# print F-values of successive temperatures:
fvalue <- c()
for(T in temp_set){
  
  tmp <- df_mature[df_mature$temperature == T,]
  mod <- aov(energy ~ walker, tmp)
  v <- fv(mod)
  fvalue <- c(fvalue, v[1])
  if(T %in% temp_subset){
    s <- sprintf("%6.4f\t%8.4f\t%.8f", T, v[1], v[2])
    cat(s, '\n')
  }
}

# plot the F-value with cooling
plot(temp_set, fvalue, 
     xlim = c(0.1,0),
     col='blue', 'b',
     xlab = 'Temperature',
     ylab = 'F-value')

