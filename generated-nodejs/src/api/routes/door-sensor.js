const Router = require('hermesjs/lib/router');
const {validateMessage} = require('../../lib/message-validator');
const router = new Router();
const doorSensorHandler = require('../handlers/door-sensor');
module.exports = router;

router.use('door-sensor', async (message, next) => {
  try {
    await validateMessage(message.payload,'door-sensor','door-sensor-data','publish');
    await doorSensorHandler.DoorSensor({message});
    next();
  } catch (e) {
    next(e);
  }
});
