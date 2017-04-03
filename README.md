# Setup

```
# Clone this repo to $GOPATH/src/the-interceptor/

glide install
```

# DB Setup

You need to install ruby to setup your database.

```
gem install convergence
cp db/database.{yml.example,yml}
mysql -u root -e 'create database the_interceptor;'

convergence -c db/database.yml -i db/schema.rb --apply
```
