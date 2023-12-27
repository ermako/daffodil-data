# daffodil-data

It's a basic read-only API that reads notes from a CSV file and returns them by ID or at random.

## Setup
Include a CSV file in the root directory called `stickies.csv`. The records in the CSV should be in the format: `id,month,year,text,` where `id` starts at zero and increases by one each record.

for example:
```csv
0,1,2024,test for january 2024,
1,3,1998,test for march 1998,
2,5,1997,test for may 1997,
```

Currently the API runs on localhost:8484.
