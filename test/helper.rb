module CLI
  extend self

  def new_canvas(cmd = 'canvas new')
    c_url = `#{cmd}`.strip
    c_url.split('/').last
  end

  def pull_canvas(id)
    `canvas pull #{id}`.strip
  end
end
