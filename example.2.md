Same input (two different examples)

```json
{"height":16,"tiles":{"0,0":{"id":1578,"tileset":"Objects"},"0,10":{"id":1506,"tileset":"Objects"},"0,8":{"id":1504,"tileset":"Objects"},"0,9":{"id":1505,"tileset":"Objects"},"3,10":{"id":1731,"tileset":"Objects"},"3,11":{"id":1732,"tileset":"Objects"},"3,12":{"id":1733,"tileset":"Objects"},"3,13":{"id":1734,"tileset":"Objects"},"3,14":{"id":1735,"tileset":"Objects"},"3,15":{"id":1736,"tileset":"Objects"},"4,1":{"id":837,"tileset":"Objects_animated"},"4,2":{"id":838,"tileset":"Objects_animated"},"4,3":{"id":839,"tileset":"Objects_animated"},"4,4":{"id":840,"tileset":"Objects_animated"},"4,5":{"id":841,"tileset":"Objects_animated"}},"width":16,"x":16,"y":0}
```

```json
{
      "height": 16,
      "tiles": {
        "0,0": {
          "id": 1578,
          "tileset": "Objects"
        },
        "0,10": {
          "id": 1506,
          "tileset": "Objects"
        },
        "0,8": {
          "id": 1504,
          "tileset": "Objects"
        },
        "0,9": {
          "id": 1505,
          "tileset": "Objects"
        },
        "3,10": {
          "id": 1731,
          "tileset": "Objects"
        },
        "3,11": {
          "id": 1732,
          "tileset": "Objects"
        },
        "3,12": {
          "id": 1733,
          "tileset": "Objects"
        },
        "3,13": {
          "id": 1734,
          "tileset": "Objects"
        },
        "3,14": {
          "id": 1735,
          "tileset": "Objects"
        },
        "3,15": {
          "id": 1736,
          "tileset": "Objects"
        },
        "4,1": {
          "id": 837,
          "tileset": "Objects_animated"
        },
        "4,2": {
          "id": 838,
          "tileset": "Objects_animated"
        },
        "4,3": {
          "id": 839,
          "tileset": "Objects_animated"
        },
        "4,4": {
          "id": 840,
          "tileset": "Objects_animated"
        },
        "4,5": {
          "id": 841,
          "tileset": "Objects_animated"
        }
      },
      "width": 16,
      "x": 16,
      "y": 0
    }
```


output for both input examples
```json
{
 "x": 16, 
 "y": 0, 
 "width": 16, 
 "height": 16, 
 "tiles": {
  "0,0" : { "tileset": "Objects"         , "id": 1578 }, 
  "0,8" : { "tileset": "Objects"         , "id": 1504 }, 
  "0,9" : { "tileset": "Objects"         , "id": 1505 }, 
  "0,10": { "tileset": "Objects"         , "id": 1506 }, 
  "3,10": { "tileset": "Objects"         , "id": 1731 }, 
  "3,11": { "tileset": "Objects"         , "id": 1732 }, 
  "3,12": { "tileset": "Objects"         , "id": 1733 }, 
  "3,13": { "tileset": "Objects"         , "id": 1734 }, 
  "3,14": { "tileset": "Objects"         , "id": 1735 }, 
  "3,15": { "tileset": "Objects"         , "id": 1736 }, 
  "4,1" : { "tileset": "Objects_animated", "id":  837 }, 
  "4,2" : { "tileset": "Objects_animated", "id":  838 }, 
  "4,3" : { "tileset": "Objects_animated", "id":  839 }, 
  "4,4" : { "tileset": "Objects_animated", "id":  840 }, 
  "4,5" : { "tileset": "Objects_animated", "id":  841 }
 }
}
```