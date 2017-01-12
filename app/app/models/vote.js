import DS from 'ember-data';

export default DS.Model.extend({
    proposal: DS.belongsTo('proposal'),
    voted: DS.attr('boolean')
});
