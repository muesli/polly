import Ember from 'ember';
import AuthenticatedRouteMixin from 'ember-simple-auth/mixins/authenticated-route-mixin';

export default Ember.Route.extend(AuthenticatedRouteMixin, {
    activate: function() {
        this._super();
        Ember.$('grantsAccepted').button('toggle');
    },

    model() {
       return Ember.RSVP.hash({
         acceptedProposals: this.store.query('proposal', {
             accepted: true,
             ended: true
         })
       });
     },

     setupController(controller, models) {
       controller.set('acceptedProposals', models.acceptedProposals);
     }
});
