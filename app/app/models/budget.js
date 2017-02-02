import DS from 'ember-data';

export default DS.Model.extend({
  available_small: DS.attr('number'),
  value: DS.attr('number'),
  maxvalue: DS.attr('number'),
  period_end: DS.attr('isodate'),
  large_grant_period_end: DS.attr('isodate')
});
