var KiloBytes = 1 << 10;
var MegaBytes = 1 << 20;
var GigaBytes = 1 << 30;

var round2decimals = function (value) {
  if (!value) return 0;
  if (value >= 100) {
    return 100;
  }
  if (value <= 0) {
    return 0;
  }
  return Math.round(value*100)/100;
};

var bytesHumanize = function (value) {
  if (value < KiloBytes) {
    return value+'B';
  }
  if (value < MegaBytes) {
    return round2decimals(value / KiloBytes)+'KiB';
  }
  if (value < GigaBytes) {
    return round2decimals(value / MegaBytes)+'MiB';
  }
  return round2decimals(value / GigaBytes)+'GiB';
};
