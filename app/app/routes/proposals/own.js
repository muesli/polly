import Ember from 'ember';
import AuthenticatedRouteMixin from 'ember-simple-auth/mixins/authenticated-route-mixin';

const { inject: { service } } = Ember;

export default Ember.Route.extend(AuthenticatedRouteMixin, {
    currentUser: service('current-user'),

    activate: function() {
        this._super();
        Ember.$('grantsOwn').button('toggle');
    },

    model() {
       return Ember.RSVP.hash({
         ownProposals: this.store.query('proposal', {
             user_id: this.get('currentUser').get('user').get('id')
         })
       });
     },

     setupController(controller, models) {
       controller.set('ownProposals', models.ownProposals);
     }
});
