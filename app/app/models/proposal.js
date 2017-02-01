import DS from 'ember-data';

export default DS.Model.extend({
  title: DS.attr('string'),
  description: DS.attr('string'),
  activities: DS.attr('string'),
  url: DS.attr('string'),
  user: DS.belongsTo('user'),
  contact: DS.attr('string'),
  recipient: DS.attr('string'),
  recipient2: DS.attr('string'),
  value: DS.attr('number'),
  granttype: DS.attr('string'),
  starts: DS.attr('isodate'),
  ends: DS.attr('isodate'),
  finished_date: DS.attr('isodate'),
  ended: DS.attr('boolean'),
  accepted: DS.attr('boolean'),
  moderated: DS.attr('boolean'),
  votes: DS.attr('number')
});
