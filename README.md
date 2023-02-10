# Jupiter

Very much a work in progress...

Jupiter is a collection of software for my skycam project. For now it'll run on a Raspberry Pi and upload photo/video to
DigitalOcean via their blob storage, Spaces (which uses AWS S3...).

## Moneta
Moneta is the name of the DigitalOcean Spaces bucket and it will hold all photos taken by the skycam as well as any
timelapse videos created.

## Minerva
A database (also on DigitalOcean) to store photo/timelapse metadata.
