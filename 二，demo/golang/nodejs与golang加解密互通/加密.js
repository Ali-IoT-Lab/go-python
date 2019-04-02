
let crypto = require("crypto")
let iv = [0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d,
  0x0e, 0x0f];



console.log(encrypt_with_aes("asdasdasdasdasd","test"))

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
  console.log('encode message: ' + encrypted);
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