import Ember from 'ember';

export default Ember.View.extend({
  didInsertElement: function () {
    this.updateDomStuff();
  },

  updateDomStuff: Ember.observer('controller.model.compiledMarkdown', function () {
    Ember.run.scheduleOnce('afterRender', this, function() {
      var el = this.get('element');
      if (!el) {
        return;
      }
      // Render tex
      Ember.$(el).find('script').each(function (i, e) {
        var t = e.getAttribute('type');
        if (t.search('math/tex') === 0) {
          katex.render(e.textContent,
                       e.parentElement,
                       {
                         displayMode: t.search('mode=display') !== -1
                       });
        }
      });

      // Fix internal links
      Ember.$(el).find('.page-markdown a').each((i, el) => {
        var href = el.getAttribute('href');
        if (href[0] === '/') {
          Ember.$(el).click((ev) => {
            var router = this.get('controller.target.router');
            router.transitionTo(href);
            ev.preventDefault();
          });
        }
      });
    });
  }),
});
