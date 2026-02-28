# LostThingsSearch service
## Current roles:
```
1 - superadmin
2 - admin
3 - institution_administrator
4 - staff
5 - teacher
6 - parent
7 - student
```

## Installing `cron` schedule:

```
sh -c "( crontab -l; cat ./crontab.tasks )" | crontab -
```
