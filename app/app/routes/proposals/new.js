import Ember from 'ember';
import AuthenticatedRouteMixin from 'ember-simple-auth/mixins/authenticated-route-mixin';

export default Ember.Route.extend(AuthenticatedRouteMixin, {
    model() {
       return Ember.RSVP.hash({
           budget: this.store.findAll('budget')
       });
     },

     setupController(controller, models) {
       controller.set('budget', models.budget);
     }
});
