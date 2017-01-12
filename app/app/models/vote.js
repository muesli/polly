import DS from 'ember-data';

export default DS.Model.extend({
//    user: DS.belongsTo('user'),
    proposal: DS.belongsTo('proposal'),
    voted: DS.attr('boolean')
});
