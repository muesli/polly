import DS from 'ember-data';

export default DS.Model.extend({
  value: DS.attr('number'),
  maxvalue: DS.attr('number'),
  period_end: DS.attr('isodate')
});
