import DS from 'ember-data';

export default DS.Model.extend({
  email: DS.attr('string'),
  about: DS.attr('string'),
  activated: DS.attr('boolean')
});
