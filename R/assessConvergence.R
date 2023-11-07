
library(readr)
library(dplyr)
library(ggplot2)

# read data 
filename <- "./data/gb_1.0.csv"
df <- read_csv(filename,
                 show_col_types = FALSE)
df['walker'] <- factor(df$walker)

# plot by temperature - but only m~12 values
temp_set <- rev(sort(unique(df$temperature)))
n <- length(temp_set)
m <- 10
temp_subset <- temp_set[seq(1,n,ceiling(n/m))]

df %>% 
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
  
  tmp <- df[df$temperature == T,]
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
     xlim = c(1,0),
     col='blue', 'b',
     xlab = 'Temperature',
     ylab = 'F-value')

