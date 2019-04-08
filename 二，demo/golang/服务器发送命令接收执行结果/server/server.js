'use strict'
const Base64 = require('js-base64').Base64;
const stringbuffer = require('stringbuffer');
const fs = require('fs');
const path = require('path');
let app = require('http').createServer(handler)
let io = require('socket.io')(app)


var sbb = new stringbuffer()
app.listen(3000)
let crypto = require("crypto")
let iv = [0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d,
  0x0e, 0x0f];

console.log('Listening at http://localhost:3000/')


function handler (req, res) {
  res.writeHead(200)
  res.end('Testing server for http://github.com/wedeploy/gosocketio example.')
}

io.on('connection', function (socket) {
  // console.log('Connecting %s.', socket.id)
 // console.log(socket.request.headers)
  //socket.emit("/message", "{10, \"main\", \"using emit\"}")

  var head = socket.request.headers
  console.log(head.ticket)
  const privateKey = fs.readFileSync(path.join(__dirname, "./private.pem"), "utf8");
  var plainText = crypto.privateDecrypt(privateKey, Buffer.from(head.ticket, "base64")).toString("utf8");
  console.log(" ----------------------client encode string-----------------------------");
  console.log(plainText);

  socket.on('messgae', (location) => {
    // fail booking 50% of the requests

    // const o = Buffer.from('o')
    // console.log(t.toString('base64'));
    // const p = Buffer.from('p')
    // console.log(p.toString('base64'));
    // const n = Buffer.from('\n')
    // console.log(n.toString('base64'));
    // console.log('client message!----------------------------------------------')
    // console.log(t.toString('base64').length)
    // console.log(o.toString('base64'))
    // console.log(p.toString('base64'))
    // console.log(n.toString('base64'))

    // socket.binaryType = 'arraybuffer'
    // var buf = new ArrayBuffer(4)
    // var bufView = new Uint8Array(buf)
    // bufView[0] = 't';
    // bufView[1] = 'o';
    // bufView[2] = 'p';
    // bufView[3] = '\n';

    //socket.send(sbb.append("top"))

    // socket.send(encrypt_with_aes("asdasdasdasdasd","t"))
    // socket.send(encrypt_with_aes("asdasdasdasdasd","o"))
    // socket.send(encrypt_with_aes("asdasdasdasdasd","p"))
    // socket.send(encrypt_with_aes("asdasdasdasdasd","\n"))

    // console.log("kkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkk")
    // console.log(location)
    //
    //
    //
    // console.log("99923456789087654323456789")
    // console.log(Base64.decode(location))

    const t = Buffer.from('testtttt');
    //console.log(t.toString('base64'));

    socket.send(encrypt_with_aes("asdasdasdasdasd",t.toString('base64')))

   // socket.send("s")
   //  socket.send("\n")


  })

  socket.on('/join', (location) => {
    // fail booking 50% of the requests
    console.log('joinjoinjoinjoinjoinjoinjoin!')
    console.log(location)

  })

  socket.on('result', (location) => {
    console.log(location)

  })


  socket.on('find_tickets', (route) => {
    console.log('Quote for tickets from %s to %s requested.', route.From, route.To)
  })

  socket.on('error', (err) => {
    console.error(err)
  })

  socket.on('disconnect', () => {
    console.log('Disconnecting %s.', socket.id)
  })
})








function encrypt_with_aes(key, message) {
  let md5 = crypto.createHash('md5').update(key).digest('hex');
  const cipher = crypto.createCipheriv(
    'aes-128-cbc',
    new Buffer(md5, 'hex'),
    new Buffer(iv)
  );
  // cipher.setAutoPadding(true);
  var encrypted = cipher.update(message, 'utf8', 'base64');
  encrypted += cipher.final('base64');
  //console.log('encode message: ' + encrypted);
  return encrypted;
}

function decrypt_with_aes(key, message) {
  let md5 = crypto.createHash('md5').update(key).digest('hex');
  const decipher = crypto.createDecipheriv(
    'aes-128-cbc',
    new Buffer(md5, 'hex'),
    new Buffer(iv)
  );
  var decrypted = decipher.update(message, 'base64', 'utf8');
  decrypted += decipher.final('utf8');
  console.log('decode message: ' + decrypted);
  return decrypted;
}