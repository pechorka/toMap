# tomap - gorm plugin that allows queering data in to a map

## How it works

For now, tomap only works with the following maps:

map[uint]ModelType

ModelType is a struct with conventional for gorm primary key ID (https://gorm.io/docs/conventions.html#ID-as-Primary-Key) and uint is an ID field value.

Before query tomap replaces map as the target value with a slice. So if you have map[uint]User it will be replaced with []User slice.

After query tomap populates the original target map with values from slice.


## Usage

tomap.RegisterCallbacks(gormDbInst)