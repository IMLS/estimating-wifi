import { mount } from '@vue/test-utils'
import USWDSContent from './USWDSContent.vue'

describe('USWDSContent', () => {
  it('should convert multiline text to separate paragraphs', () => {
    const multilineContent = 'Line 1 \n Line 2';
    const wrapper = mount(USWDSContent, { props: { multilineContent } })

    expect(wrapper.findAll('p')).toHaveLength(2)
  })
})