module CLI
  extend self

  def bin
    './canvas'
  end

	def delete(id)
		`#{CLI.bin} delete #{id}`
	end

  def new_canvas(cmd = nil)
    cmd ||= "#{CLI.bin} new"
    c_url = `#{cmd}`.strip
    c_url.split('/').last
  end

  def pull_canvas(id)
    `./#{self.bin} pull #{id}`.strip
  end
end
