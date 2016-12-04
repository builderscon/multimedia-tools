Video Splitter
==============

If you are not contracting out video editing, you are going to need to
split out talk videos from the original unedited raw (long) footage.

This simple Perl script does exactly that, and a little bit more.
This script requires `ffmpeg`

Usage
=====

```
perl ./video-split.pl <source file> <config file>
```

Config File
===========

The config file is in JSON format, and contains an array of objects. 

## Parameters

### filename

The name of the output file.

### start

Timestamp since the beginning of the source movie to start extracting contents
out from. Format should be in `hh:mm:ss`

### end

Same as start, but denotes the end of the extraction. Format should be in `hh:mm:ss`

### cropblur

Specifies a section of the video that should be blurred out. This often happens
when the speaker has material that can be shared for an audience but not online, or maybe there was a single slide that *just* needed to be edited out.

The cropblur attribute is also an object. Here are the possible keys:

| key     | required | notes |
|---------|----------|-------|
| x       | YES      | x offset of topleft corner to start the blur |
| y       | YES      | y offset of topleft corner to start the blur |
| width   | YES      | width of the rectangle to blur |
| height  | YES      | height of the rectangle to blur |
| between | NO       | an array containing "start" and "end" timestamps. This can control when in the timeline the blur is applied to |

The cropblur currently only supports one specification per output file, but
this probably needs to be changed in the future. Please submit a PR.

Here's a sample, used in the very first [builderscon tokyo 2016](https://builderscon.tokyo/builderscon/tokyo/2016):

```json
[
  { "filename": "A_1000_lestrrat.mov", "start": "00:46:38", "end": "00:54:09" },
  { "filename": "A_1010_mattn.mov", "start": "00:56:35", "end": "01:48:44" },
  { "filename": "A_1120_uzulla.mov", "start": "02:06:41", "end": "03:19:00" },
  { "filename": "A_1310_miyake_youichiro.mov", "start": "03:56:54", "end": "04:57:44",
    "cropblur": {
      "width": 840,
      "height": 640,
      "x": 120,
      "y": 75,
      "between": [
        { "start": "00:05:10", "end": "00:05:52" },
        { "start": "00:07:28", "end": "00:07:57" },
        { "start": "00:24:28", "end": "00:25:06" },
        { "start": "00:25:11", "end": "00:26:59" },
        { "start": "00:39:36", "end": "00:39:55" },
        { "start": "00:39:59", "end": "00:40:17" },
        { "start": "00:40:22", "end": "00:40:41" },
        { "start": "00:58:59", "end": "00:59:15" },
        { "start": "00:59:20", "end": "01:00:35" }
      ]
    }
  },
  { "filename": "A_1420_mumoshu.mov", "start": "05:06:40", "end": "06:07:11" },
  { "filename": "A_1530_kazsato.mov", "start": "06:16:37", "end": "06:46:31" },
  { "filename": "A_1600_shibayu.mov", "start": "06:48:09", "end": "07:18:32" },
  { "filename": "A_1640_cho45.mov", "start": "07:27:44", "end": "08:26:48" }
]
```