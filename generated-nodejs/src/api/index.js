const Hermes = require('hermesjs');
const app = new Hermes();
const { cyan, gray, yellow } = require('colors/safe');
const buffer2string = require('./middlewares/buffer2string');
const string2json = require('./middlewares/string2json');
const json2string = require('./middlewares/json2string');
const logger = require('./middlewares/logger');
const errorLogger = require('./middlewares/error-logger');
const config = require('../lib/config');
const serverConfig = config.broker.kafka;
const KafkaAdapter = require('hermesjs-kafka');
const doorSensor = require('./routes/door-sensor.js');
const fs = require('fs')

serverConfig.ssl.ca = fs.readFileSync('ca.pem');
serverConfig.ssl.key = fs.readFileSync('service.key');
serverConfig.ssl.cert = fs.readFileSync('service.cert');

app.addAdapter(KafkaAdapter, serverConfig);

app.use(buffer2string);
app.use(string2json);
app.use(logger);

// Channels
console.log(cyan.bold.inverse(' SUB '), gray('Subscribed to'), yellow('door-sensor'));
app.use(doorSensor);

app.use(errorLogger);
app.useOutbound(errorLogger);
app.useOutbound(logger);
app.useOutbound(json2string);

app
  .listen()
  .then((adapters) => {
    console.log(cyan.underline(`${config.app.name} ${config.app.version}`), gray('is ready!'), '\n');
    adapters.forEach(adapter => {
      console.log('🔗 ', adapter.name(), gray('is connected!'));
    });
  })
  .catch(console.error);
