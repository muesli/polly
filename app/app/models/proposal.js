import DS from 'ember-data';

export default DS.Model.extend({
  title: DS.attr('string'),
  description: DS.attr('string'),
  url: DS.attr('string'),
  recipient: DS.attr('string'),
  value: DS.attr('number'),
  granttype: DS.attr('string'),
  ends: DS.attr('isodate'),
  ended: DS.attr('boolean'),
  moderated: DS.attr('boolean'),
  votes: DS.attr('number')
});
