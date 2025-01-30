# Golang Presidential-Action Tracker

This bot will track the 10 latest presidential actions listed on whitehouse.gov.

# Install

```
git clone https://github.com/NoahSoto/executive-order-tracker
cd ./executive-order-tracker
pip install -r requirements.txt
go run main.go &
```


# Usage

Initialize the bot with exeuctive orders, and start the golang worker to check for new actions by running

```
!orders
```
Then once thats finished you can view previous orders by running 

```
!ls
```

And finally you can retrieve the title, URL, and a brief NLTK based synopsis via

```
!view [1-10]
```

