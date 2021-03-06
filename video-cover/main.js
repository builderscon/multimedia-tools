(function() {
// var conferenceId = 'fdd2e7aa-b234-4aac-bfb1-e58cb334e8a7'
var images = [];

var baseFontSize = 60;
var coverWidth = 1280;
var coverHeight = 720;
var fontName = "'Mplus 1p'";
var fonts = [
  "'Mplus 1p'",
  "'Rounded Mplus 1c'",
  "Hannari",
  "Kokoro",
  "'Sawarabi Mincho'",
  "'Sawarabi Gothic'",
  "Nikukyu",
  "'Nico Moji'",
  "'Noto Sans Japanese'",
];
var metaKeys = ['title', 'subTitle'];
var sessionKeys = ['image', 'title', 'name', 'fontSize'];
var storageKey = 'data';
var previewEnabled = false;

// Canvas

function drawBackground(ctx, img) {
  if (!img) {
    console.log('img is required');
    return;
  }

  var scale = ctx.canvas.width/img.originalWidth;
  ctx.save()
  ctx.setTransform(scale, 0, 0, scale, 0, 0);
  ctx.drawImage(img, 0, 0);
  ctx.restore()
}

function drawWatermark(ctx, img) {
  // hardcoded 200x200 image
  var watermark = document.getElementById('watermark');
  var scale = ctx.canvas.width/img.originalWidth;
  ctx.save();
  ctx.setTransform(scale, 0, 0, scale, 0, 0);
  // This is weird shit, but the height of the canvas != image height,
  // so in order to calculate 205px above bottom right corner, we need
  // to scale *BACK* the canvas height then subtract 205. This result is
  // then converted via the above setTransform into the appropriate canvas size
  ctx.drawImage(watermark, img.originalWidth - 205, (ctx.canvas.height / scale) - 205);
  ctx.restore()
}

function drawShadow(ctx) {
  var grd = ctx.createLinearGradient(coverWidth/4, 0, coverWidth/2, 0);
  grd.addColorStop(0, 'rgba(0,0,0,1)');
  grd.addColorStop(1, 'rgba(0,0,0,0)');
  ctx.fillStyle = grd;
  ctx.fillRect(0, 0, coverWidth, coverHeight);
}

function drawTitle(ctx, title, fontSize) {
  if (!title) {
    return;
  }

  ctx.fillStyle = '#fff';
  var text = '';
  var lines = title.split("\\n");
  ctx.font = 'bold ' + fontSize + 'px ' + fontName;
  for (var i = 0; i < lines.length; i++) {
    ctx.fillText(lines[i], 10, 120 + ((fontSize - baseFontSize + 20) *1.5 + 50) * i, coverWidth/2);
  }
}

function drawMeta(ctx, name, conference, date) {
  ctx.font = "bold 60px " + fontName;
  ctx.fillText(name || 'NAME', 15, coverHeight - 150);

  ctx.font = "bold 40px " + fontName;
  ctx.fillText(conference || 'builderscon tokyo 2016', 15, coverHeight - 80);
  ctx.fillText(date || 'Dec 3, 2016', 15, coverHeight - 35);
}

function doCanvas(canvas, img, title, name, fontSize, conference, date) {
  var ctx = canvas.getContext('2d');
  drawBackground(ctx, img);

  var scale = canvas.width/coverWidth;
  ctx.scale(scale, scale);
  drawShadow(ctx);
  drawTitle(ctx, title, fontSize);
  drawMeta(ctx, name, conference, date);
  drawWatermark(ctx, img);
}

// Event

function onFetch() {
  var inp = document.querySelector('input#conference_id')
  fetchConferenceData(inp.value)
    .then(_ => buildTable())
    .catch(err => console.log(err));
}

function onChangeFont(ev) {
  fontName = fonts[ev.target.selectedIndex];
  document.body.style.fontFamily = fontName;
  var thumb = document.getElementById('settings_thumbnail');
  if (thumb.classList.contains('is-visible')) {
    setTimeout(_ => {
      buildThumbnails();
    }, 800);
  }
}

function onChangeInput(ev) {
  var collectValues = () => {
    var div = document.getElementById('settings_table');
    var table = div.querySelector('table');
    var sessions = [];
    table.querySelectorAll('.session').forEach(s => {
      var session = {}
      s.querySelectorAll('input').forEach(inp => {
        session[inp.name] = inp.value;
      });
      sessions.push(session);
    });

    var data = {
      meta: {},
      sessions: sessions
    }

    metaKeys.forEach(key => {
      var inp = div.querySelector('input[name=' + key + ']');
      data.meta[key] = inp.value;
    });

    return data
  }

  // update localStorage using table values
  var data = collectValues()
  localStorage.setItem(storageKey, JSON.stringify(data));

  // preview cover
  if (!previewEnabled) {
    return
  }
  var row = ev.target.parentNode.parentNode;
  if (!row.classList.contains('session')) {
    return
  }

  var table = ev.target.parentNode.parentNode.parentNode;
  var preview = document.getElementById('cover_preview');
  table.insertBefore(preview, row.nextSibling);
  preview.classList.add('is-preview');

  // draw preview canvas
  var canvas = preview.querySelector('canvas');
  var s = {}
  row.querySelectorAll('input').forEach(inp => {
    s[inp.name] = inp.value;
  });
  var img = images.find(a => a.src === s.image);
  if (img) {
    doCanvas(canvas, img, s.title, s.name, s.fontSize, data.conference, data.date);
  } else {
    fetchImage(s.image, img => {
      doCanvas(canvas, img, s.title, s.name, s.fontSize, data.conference, data.date);
    });
  }
}

function onShowTable() {
  document.getElementById('settings_thumbnail').classList.remove('is-visible');
  buildTable();
}

function onShowThumbnail() {
  document.getElementById('settings_table').classList.remove('is-visible');
  buildThumbnails();
}

function onDownload() {
  var str = localStorage.getItem(storageKey);
  if (!str) {
    console.log('No data');
    return;
  }

  var zip = new JSZip();
  var imgzip = zip.folder('images');
  var data = JSON.parse(str);
  var count = 0;

  data.sessions.forEach(a => {
    var img = images.find(i => i.src === a.image)
    if (!img) {
      return;
    }

    var canvas = document.createElement('canvas');
    canvas.width = coverWidth;
    canvas.height = coverHeight;
    doCanvas(canvas, img, a.title, a.name, a.fontSize, data.meta.title,  data.meta.subTitle);
    var imgData = canvas.toDataURL('image/jpg')
    imgzip.file(a.title + '.jpg', imgData.split(',')[1], {base64: true});
  });

  imgzip.generateAsync({type:"blob"})
    .then(blob => {
      saveAs(blob, "images.zip");
    })
    .catch(err => {
      console.log('zip error', err);
    });
}

// Builder

function buildTable() {
  var str = localStorage.getItem(storageKey);
  if (!str) {
    console.log('No data')
    return
  }

  var data = JSON.parse(str);

  var div = document.getElementById('settings_table')
  div.classList.add('is-visible');
  div.innerHTML = '';

  var h3 = document.createElement('h3')
  h3.innerHTML = 'Conference';
  div.appendChild(h3);
  metaKeys.forEach(key => {
    var label = document.createElement('label');
    label.innerHTML = key;
    label.htmlFor = key;
    div.appendChild(label);
    var inp = document.createElement('input');
    inp.name = key;
    inp.value = data.meta[key]
    div.appendChild(inp);
  });

  var h4 = document.createElement('h4');
  h4.innerHTML = 'Sessions';
  div.appendChild(h4);
  // To enable preview
  var chbox = document.createElement('input');
  chbox.id = 'enable_preview_check_box';
  chbox.type = 'checkbox';
  chbox.checked = previewEnabled;
  chbox.addEventListener('change', ev =>{
    previewEnabled = ev.target.checked
    if (!previewEnabled) {
      document.getElementById('cover_preview').remove('is-preview');
    }
  });
  div.appendChild(chbox);
  var label = document.createElement('label');
  label.innerHTML = 'Enable preview';
  label.htmlFor = chbox.id;
  div.appendChild(label);

  var table = document.createElement('table');
  div.appendChild(table);
  // table header
  var tr = document.createElement('tr');
  table.appendChild(tr);
  sessionKeys.forEach(a => {
    var th = document.createElement('th');
    th.innerHTML = a;
    tr.appendChild(th);
  })
  // cover preview
  var tr = document.createElement('tr');
  tr.id = 'cover_preview';
  if (previewEnabled) {
    tr.classList.add('is-preview');
  }
  table.appendChild(tr);
  var td = document.createElement('td');
  td.colSpan = sessionKeys.length;
  tr.appendChild(td);
  var canvas = document.createElement('canvas');
  canvas.width = coverWidth;
  canvas.height = coverHeight;
  td.appendChild(canvas);

  data.sessions.forEach(a => {
    var tr = document.createElement('tr');
    tr.classList.add('session');
    table.appendChild(tr);

    sessionKeys.forEach(key => {
      var td = document.createElement('td');
      tr.appendChild(td);
      var inp = document.createElement('input');
      inp.name = key;
      if (a[key]) {
        inp.value = a[key];
      }
      td.appendChild(inp);
    });
  });

  div.querySelectorAll('input')
    .forEach(a => a.addEventListener('input', onChangeInput))
}

function buildThumbnails() {
  var str = localStorage.getItem(storageKey);
  if (!str) {
    console.log('No data')
    return
  }

  var data = JSON.parse(str);

  var div = document.getElementById('settings_thumbnail')
  div.classList.add('is-visible');
  div.innerHTML = '';
  var ul = document.createElement('ul');
  div.appendChild(ul);

  images = [];
  data.sessions.forEach(a => {
    var canvas = document.createElement('canvas');
    var img = fetchImage(a.image, img => {
      doCanvas(canvas, img, a.title, a.name, a.fontSize, data.meta.title, data.meta.subTitle);
    })
    img.addEventListener('click', event => {
      // TODO
    });
    images.push(img);

    var li = document.createElement('li');
    canvas.width = 640;
    canvas.height = 360;
    li.appendChild(canvas);
    ul.appendChild(li);
  })
}

// Other

function fetchConferenceData(id) {
  var url = 'https://api.builderscon.io/v2/conference/lookup?id=' + id + '&lang=ja';
  return fetch(url)
    .then(data => data.json())
    .then(conf => {
      var url = 'https://api.builderscon.io/v2/session/list?conference_id=' + id + '&lang=ja';
      return fetch(url)
        .then(data => data.json())
        .then(json => {
          var sessions = json
            .filter(a => a.video_url)
            .map(a => {
              var name = a.speaker.nickname;
              if (a.speaker.first_name && a.speaker.last_name && a.speaker.first_name != 'Unknown' && a.speaker.last_name != 'Unknown') {
                name = a.speaker.first_name + " " + a.speaker.last_name;
              }
              return {
                title: a.title,
                name: name,
                fontSize: baseFontSize
              }
            })
  
          var data = {
            meta: {
              title: conf.title,
              subTitle: conf.sub_title
            },
            sessions: sessions 
          }
          localStorage.setItem(storageKey, JSON.stringify(data));
          return data;
        })
    })
}

function fetchImage(src, callback) {
  var img = new Image();
  img.onload = function() {
    img.originalWidth = img.width;
    callback(img);
  }
  img.setAttribute('crossOrigin', 'anonymous');
  img.src = src;
  return img;
}

// main
var sel = document.getElementById('fonts');
fonts.forEach(a => {
  var opt = document.createElement('option');
  opt.innerHTML = a;
  sel.append(opt);
});
sel.addEventListener('change', onChangeFont);
document.getElementById('on_show_table').addEventListener('click', onShowTable);
document.getElementById('on_show_thumbnail').addEventListener('click', onShowThumbnail);
document.getElementById('on_fetch').addEventListener('click', onFetch);
document.getElementById('on_download').addEventListener('click', onDownload);

buildTable();
})();
